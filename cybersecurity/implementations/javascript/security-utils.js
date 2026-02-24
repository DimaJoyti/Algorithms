const crypto = require('crypto');

class PasswordHasher {
    static async hash(password, rounds = 12) {
        return new Promise((resolve, reject) => {
            const salt = crypto.randomBytes(16).toString('hex');
            crypto.pbkdf2(password, salt, rounds * 1000, 64, 'sha512', (err, derivedKey) => {
                if (err) reject(err);
                resolve(`pbkdf2:${rounds}:${salt}:${derivedKey.toString('hex')}`);
            });
        });
    }
    
    static async verify(password, hashedPassword) {
        const [algorithm, rounds, salt, hash] = hashedPassword.split(':');
        return new Promise((resolve, reject) => {
            crypto.pbkdf2(password, salt, parseInt(rounds) * 1000, 64, 'sha512', (err, derivedKey) => {
                if (err) reject(err);
                resolve(crypto.timingSafeEqual(
                    Buffer.from(hash, 'hex'),
                    Buffer.from(derivedKey.toString('hex'), 'hex')
                ));
            });
        });
    }
}

class JWTHandler {
    constructor(secret, algorithm = 'HS256') {
        this.secret = secret;
        this.algorithm = algorithm;
    }
    
    sign(payload, expiresIn = '1h') {
        const header = { alg: this.algorithm, typ: 'JWT' };
        const now = Math.floor(Date.now() / 1000);
        
        payload.iat = now;
        payload.exp = now + this._parseExpiry(expiresIn);
        payload.jti = crypto.randomBytes(16).toString('hex');
        
        const encodedHeader = this._base64UrlEncode(JSON.stringify(header));
        const encodedPayload = this._base64UrlEncode(JSON.stringify(payload));
        const signature = this._sign(`${encodedHeader}.${encodedPayload}`);
        
        return `${encodedHeader}.${encodedPayload}.${signature}`;
    }
    
    verify(token) {
        const [header, payload, signature] = token.split('.');
        
        const expectedSignature = this._sign(`${header}.${payload}`);
        if (!crypto.timingSafeEqual(
            Buffer.from(signature),
            Buffer.from(expectedSignature)
        )) {
            throw new Error('Invalid signature');
        }
        
        const decoded = JSON.parse(this._base64UrlDecode(payload));
        
        if (decoded.exp < Math.floor(Date.now() / 1000)) {
            throw new Error('Token expired');
        }
        
        return decoded;
    }
    
    _sign(data) {
        return crypto.createHmac('sha256', this.secret)
            .update(data)
            .digest('base64')
            .replace(/\+/g, '-')
            .replace(/\//g, '_')
            .replace(/=/g, '');
    }
    
    _base64UrlEncode(str) {
        return Buffer.from(str)
            .toString('base64')
            .replace(/\+/g, '-')
            .replace(/\//g, '_')
            .replace(/=/g, '');
    }
    
    _base64UrlDecode(str) {
        str = str.replace(/-/g, '+').replace(/_/g, '/');
        while (str.length % 4) str += '=';
        return Buffer.from(str, 'base64').toString();
    }
    
    _parseExpiry(expiresIn) {
        const match = expiresIn.match(/^(\d+)([smhd])$/);
        if (!match) throw new Error('Invalid expiry format');
        
        const [, value, unit] = match;
        const multipliers = { s: 1, m: 60, h: 3600, d: 86400 };
        return parseInt(value) * multipliers[unit];
    }
}

class RateLimiter {
    constructor(options = {}) {
        this.windowMs = options.windowMs || 60000;
        this.maxRequests = options.maxRequests || 100;
        this.requests = new Map();
        this.cleanupInterval = setInterval(() => this._cleanup(), this.windowMs);
    }
    
    middleware() {
        return (req, res, next) => {
            const key = this._getKey(req);
            const result = this._check(key);
            
            res.setHeader('X-RateLimit-Limit', this.maxRequests);
            res.setHeader('X-RateLimit-Remaining', Math.max(0, this.maxRequests - result.count));
            res.setHeader('X-RateLimit-Reset', result.resetTime);
            
            if (result.limited) {
                res.setHeader('Retry-After', Math.ceil(this.windowMs / 1000));
                return res.status(429).json({ error: 'Too many requests' });
            }
            
            next();
        };
    }
    
    _getKey(req) {
        return req.ip || req.connection.remoteAddress;
    }
    
    _check(key) {
        const now = Date.now();
        const windowStart = now - this.windowMs;
        
        if (!this.requests.has(key)) {
            this.requests.set(key, []);
        }
        
        const timestamps = this.requests.get(key);
        const validTimestamps = timestamps.filter(t => t > windowStart);
        validTimestamps.push(now);
        this.requests.set(key, validTimestamps);
        
        return {
            count: validTimestamps.length,
            limited: validTimestamps.length > this.maxRequests,
            resetTime: new Date(now + this.windowMs).toISOString()
        };
    }
    
    _cleanup() {
        const windowStart = Date.now() - this.windowMs;
        for (const [key, timestamps] of this.requests.entries()) {
            const valid = timestamps.filter(t => t > windowStart);
            if (valid.length === 0) {
                this.requests.delete(key);
            } else {
                this.requests.set(key, valid);
            }
        }
    }
}

class InputValidator {
    static sanitize(str) {
        if (typeof str !== 'string') return '';
        return str
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#x27;')
            .replace(/\//g, '&#x2F;');
    }
    
    static escapeHtml(str) {
        return this.sanitize(str);
    }
    
    static escapeJavaScript(str) {
        return str
            .replace(/\\/g, '\\\\')
            .replace(/'/g, "\\'")
            .replace(/"/g, '\\"')
            .replace(/\//g, '\\/')
            .replace(/\n/g, '\\n')
            .replace(/\r/g, '\\r')
            .replace(/\t/g, '\\t');
    }
    
    static isEmail(email) {
        return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
    }
    
    static isUrl(url) {
        try {
            new URL(url);
            return true;
        } catch {
            return false;
        }
    }
    
    static isAlpha(str) {
        return /^[a-zA-Z]+$/.test(str);
    }
    
    static isAlphanumeric(str) {
        return /^[a-zA-Z0-9]+$/.test(str);
    }
    
    static isLength(str, min, max = Infinity) {
        const len = str.length;
        return len >= min && len <= max;
    }
    
    static stripTags(str) {
        return str.replace(/<[^>]*>/g, '');
    }
    
    static allowedTags(str, tags = []) {
        const pattern = new RegExp(`<(?!\\/?(${tags.join('|')})\\b)[^>]+>`, 'gi');
        return str.replace(pattern, '');
    }
}

class XSSProtection {
    static sanitizeHTML(html, options = {}) {
        const allowedTags = options.allowedTags || ['b', 'i', 'u', 'strong', 'em', 'p', 'br'];
        const allowedAttributes = options.allowedAttributes || {};
        
        let result = html;
        
        result = result.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '');
        result = result.replace(/on\w+\s*=\s*["'][^"']*["']/gi, '');
        result = result.replace(/javascript:/gi, '');
        result = result.replace(/vbscript:/gi, '');
        result = result.replace(/data:/gi, '');
        
        result = result.replace(/<(\w+)([^>]*)>/g, (match, tag, attrs) => {
            if (!allowedTags.includes(tag.toLowerCase())) {
                return '';
            }
            
            if (allowedAttributes[tag.toLowerCase()]) {
                const cleanAttrs = attrs.replace(/on\w+\s*=/gi, '');
                return `<${tag}${cleanAttrs}>`;
            }
            
            return `<${tag}>`;
        });
        
        return result;
    }
    
    static contentSecurityPolicy(options = {}) {
        const directives = {
            'default-src': ["'self'"],
            'script-src': ["'self'"],
            'style-src': ["'self'", "'unsafe-inline'"],
            'img-src': ["'self'", 'data:'],
            'font-src': ["'self'"],
            'connect-src': ["'self'"],
            'frame-ancestors': ["'none'"],
            ...options
        };
        
        return Object.entries(directives)
            .map(([key, values]) => `${key} ${values.join(' ')}`)
            .join('; ');
    }
}

class CSRFProtection {
    constructor(secret) {
        this.secret = secret;
    }
    
    generateToken() {
        return crypto.randomBytes(32).toString('hex');
    }
    
    generateSignedToken(sessionId) {
        const token = this.generateToken();
        const signature = crypto
            .createHmac('sha256', this.secret)
            .update(`${sessionId}:${token}`)
            .digest('hex');
        return `${token}:${signature}`;
    }
    
    verifyToken(sessionId, token) {
        const [actualToken, signature] = token.split(':');
        const expectedSignature = crypto
            .createHmac('sha256', this.secret)
            .update(`${sessionId}:${actualToken}`)
            .digest('hex');
        
        return crypto.timingSafeEqual(
            Buffer.from(signature, 'hex'),
            Buffer.from(expectedSignature, 'hex')
        );
    }
    
    middleware() {
        return (req, res, next) => {
            if (['GET', 'HEAD', 'OPTIONS'].includes(req.method)) {
                return next();
            }
            
            const token = req.body._csrf || req.headers['x-csrf-token'];
            const sessionId = req.session?.id;
            
            if (!token || !sessionId || !this.verifyToken(sessionId, token)) {
                return res.status(403).json({ error: 'Invalid CSRF token' });
            }
            
            next();
        };
    }
}

class SecurityHeaders {
    static middleware() {
        return (req, res, next) => {
            res.setHeader('X-Content-Type-Options', 'nosniff');
            res.setHeader('X-Frame-Options', 'DENY');
            res.setHeader('X-XSS-Protection', '1; mode=block');
            res.setHeader('Strict-Transport-Security', 'max-age=31536000; includeSubDomains');
            res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
            res.setHeader('Permissions-Policy', 'geolocation=(), microphone=(), camera=()');
            
            res.removeHeader('X-Powered-By');
            res.removeHeader('Server');
            
            next();
        };
    }
}

module.exports = {
    PasswordHasher,
    JWTHandler,
    RateLimiter,
    InputValidator,
    XSSProtection,
    CSRFProtection,
    SecurityHeaders
};
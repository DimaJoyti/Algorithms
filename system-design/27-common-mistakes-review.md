# ❌ Common Mistakes Review

## 🎯 Overview of Common Pitfalls

System design interviews have predictable failure patterns. Understanding these common mistakes is crucial for avoiding them and demonstrating senior-level thinking.

### Categories of Mistakes
```
🚫 Communication Mistakes:
- Poor requirement clarification
- Jumping to solutions too quickly
- Not explaining thought process
- Ignoring interviewer feedback
- Poor time management

🚫 Technical Mistakes:
- Over-engineering for unrealistic scale
- Ignoring data consistency
- Forgetting about failure scenarios
- Missing security considerations
- Not considering operational complexity

🚫 Design Mistakes:
- Starting with implementation details
- Creating overly complex initial design
- Not validating against requirements
- Ignoring cost implications
- Failing to identify bottlenecks
```

## 🗣️ Communication Mistakes

### ❌ Mistake #1: Jumping to Solution Without Clarification

**What It Looks Like:**
```
Interviewer: "Design a chat system"
Bad Candidate: "Okay, so we'll use WebSockets for real-time messaging, 
               store messages in MongoDB, and use Redis for caching..."
```

**Why It's Wrong:**
- Assumes requirements without validation
- Misses opportunity to show systematic thinking
- May solve the wrong problem
- Demonstrates poor collaboration skills

**✅ Correct Approach:**
```
Good Candidate: "Great! Let me start by understanding the requirements.
                 Are we building one-on-one messaging, group chats, or both?
                 What's our expected scale in terms of users and messages?
                 Do we need features like file sharing or message history?
                 Are there any specific latency or availability requirements?"
```

**Key Lesson:** Always clarify requirements first. This shows systematic thinking and ensures you're solving the right problem.

### ❌ Mistake #2: Not Thinking Out Loud

**What It Looks Like:**
```
Candidate: [Draws diagram silently for 5 minutes]
Interviewer: "Can you explain what you're thinking?"
Candidate: "Oh, I'm just designing the database schema..."
```

**Why It's Wrong:**
- Interviewer can't follow your thought process
- Misses opportunities for guidance and feedback
- Creates awkward silences
- Doesn't demonstrate communication skills

**✅ Correct Approach:**
```
Good Candidate: "I'm thinking about the database design now. For a chat system,
                 I need to store users, conversations, and messages. Let me think
                 about the relationships... A user can be in multiple conversations,
                 and each conversation can have multiple messages. So I'm considering
                 a schema with users table, conversations table, and messages table..."
```

**Key Lesson:** Verbalize your thinking process continuously. The interviewer wants to understand how you approach problems.

### ❌ Mistake #3: Ignoring Interviewer Hints and Feedback

**What It Looks Like:**
```
Interviewer: "How would you handle the case where your database becomes a bottleneck?"
Candidate: "We could add more memory to the server."
Interviewer: "What about horizontal scaling approaches?"
Candidate: "We could get a bigger server with more CPUs."
```

**Why It's Wrong:**
- Shows inability to take feedback
- Misses learning opportunities
- Demonstrates poor collaboration
- May indicate stubbornness or lack of adaptability

**✅ Correct Approach:**
```
Good Candidate: "You're right, vertical scaling has limits. For horizontal scaling,
                 I could consider database sharding - perhaps sharding by user ID
                 or conversation ID. I could also implement read replicas to distribute
                 read load. What approach would you recommend exploring further?"
```

**Key Lesson:** Listen carefully to interviewer feedback and incorporate it into your design. This shows adaptability and collaboration.

### ❌ Mistake #4: Poor Time Management

**What It Looks Like:**
```
Time Allocation:
- Requirements: 20 minutes (too long)
- High-level design: 5 minutes (too short)
- Deep dive: 15 minutes (rushed)
- Scaling: 5 minutes (insufficient)
```

**Why It's Wrong:**
- Doesn't cover all required topics
- Shows poor prioritization
- May not demonstrate technical depth
- Creates stress and rushed explanations

**✅ Correct Approach:**
```
Good Time Allocation (45-minute interview):
- Requirements: 5-8 minutes
- High-level design: 10-15 minutes
- Deep dive: 15-20 minutes
- Scaling and optimization: 8-12 minutes
- Q&A: 3-5 minutes
```

**Key Lesson:** Practice time management and keep track of time throughout the interview.

## 🔧 Technical Mistakes

### ❌ Mistake #5: Over-Engineering for Unrealistic Scale

**What It Looks Like:**
```
Interviewer: "Design a URL shortener for a startup"
Bad Candidate: "We'll need 1000 microservices, a complex event sourcing system,
               distributed across 50 data centers, with custom consensus algorithms..."
```

**Why It's Wrong:**
- Adds unnecessary complexity
- Shows poor judgment about appropriate scale
- Increases cost and operational burden
- Demonstrates lack of practical experience

**✅ Correct Approach:**
```
Good Candidate: "For a startup URL shortener, I'll start simple. A web server,
                 a database, and a cache should handle initial scale. As we grow,
                 we can add load balancers, database replicas, and CDN.
                 Let me design for the stated requirements first."
```

**Key Lesson:** Design for the given scale, not infinite scale. Start simple and add complexity as needed.

### ❌ Mistake #6: Ignoring Data Consistency and ACID Properties

**What It Looks Like:**
```
Candidate: "We'll use eventual consistency everywhere to maximize performance."
Interviewer: "What about financial transactions in the system?"
Candidate: "Those can be eventually consistent too."
```

**Why It's Wrong:**
- Misunderstands when strong consistency is required
- Could lead to data corruption or financial losses
- Shows lack of understanding of ACID properties
- Demonstrates poor judgment about trade-offs

**✅ Correct Approach:**
```
Good Candidate: "For financial transactions, we need strong consistency and ACID
                 properties. I'll use a traditional RDBMS with transactions for
                 payment processing. For less critical data like user preferences,
                 eventual consistency might be acceptable."
```

**Key Lesson:** Understand when to use strong vs. eventual consistency based on business requirements.

### ❌ Mistake #7: Forgetting About Failure Scenarios

**What It Looks Like:**
```
Candidate: [Presents perfect happy-path design]
Interviewer: "What happens if your database goes down?"
Candidate: "Um... the system would stop working?"
```

**Why It's Wrong:**
- Shows lack of production experience
- Doesn't consider reliability requirements
- Missing critical system design thinking
- Could lead to poor system availability

**✅ Correct Approach:**
```
Good Candidate: "For database failures, I'd implement:
                 - Database replication with automatic failover
                 - Circuit breakers to handle downstream failures
                 - Graceful degradation where possible
                 - Health checks and monitoring
                 - Backup and recovery procedures"
```

**Key Lesson:** Always consider failure scenarios and design for resilience.

### ❌ Mistake #8: Missing Security Considerations

**What It Looks Like:**
```
Candidate: [Designs entire system without mentioning security]
Interviewer: "How do you handle authentication and authorization?"
Candidate: "Oh right, we should probably add some security..."
```

**Why It's Wrong:**
- Security should be designed in, not bolted on
- Shows lack of production awareness
- Could lead to serious vulnerabilities
- Demonstrates incomplete system thinking

**✅ Correct Approach:**
```
Good Candidate: "For security, I'll implement:
                 - JWT-based authentication with proper token management
                 - Role-based authorization for different user types
                 - HTTPS everywhere with proper certificate management
                 - Input validation and sanitization
                 - Rate limiting to prevent abuse
                 - Audit logging for security events"
```

**Key Lesson:** Consider security throughout the design, not as an afterthought.

## 🏗️ Design Mistakes

### ❌ Mistake #9: Starting with Implementation Details

**What It Looks Like:**
```
Candidate: "First, let me design the database schema in detail...
           [Spends 15 minutes on table structures]
           Now let me think about the API endpoints...
           [Spends 10 minutes on REST API details]"
```

**Why It's Wrong:**
- Gets lost in low-level details too early
- Doesn't establish overall architecture
- May run out of time for important topics
- Shows poor prioritization

**✅ Correct Approach:**
```
Good Candidate: "Let me start with the high-level architecture:
                 [Draws boxes for major components]
                 We have clients, load balancers, web servers, application servers,
                 databases, and caches. Let me explain the data flow...
                 [Later] Now let me dive deeper into the database design..."
```

**Key Lesson:** Start with high-level architecture, then progressively add detail.

### ❌ Mistake #10: Creating Overly Complex Initial Design

**What It Looks Like:**
```
Initial Design:
[Complex diagram with 20+ components, multiple message queues,
 microservices, event sourcing, CQRS, multiple databases]
```

**Why It's Wrong:**
- Difficult to explain and understand
- Hard to validate against requirements
- Premature optimization
- Shows poor judgment about complexity

**✅ Correct Approach:**
```
Initial Design:
[Simple diagram with 4-5 major components]
"This is my initial design. As we discuss scaling challenges,
 I can add complexity where needed."
```

**Key Lesson:** Start simple and add complexity incrementally based on requirements.

### ❌ Mistake #11: Not Validating Design Against Requirements

**What It Looks Like:**
```
Candidate: [Presents complete design]
Interviewer: "How does this handle the requirement for real-time notifications?"
Candidate: "Oh, I forgot about that requirement..."
```

**Why It's Wrong:**
- Design doesn't meet stated requirements
- Shows poor attention to detail
- Demonstrates lack of systematic validation
- May require significant rework

**✅ Correct Approach:**
```
Good Candidate: "Let me validate this design against our requirements:
                 ✓ Handles 1M users - yes, with load balancing and sharding
                 ✓ Sub-100ms latency - yes, with caching and CDN
                 ✓ Real-time updates - yes, with WebSocket connections
                 ✓ 99.9% availability - yes, with redundancy and failover"
```

**Key Lesson:** Regularly validate your design against the stated requirements.

### ❌ Mistake #12: Ignoring Cost Implications

**What It Looks Like:**
```
Candidate: "We'll replicate all data to 10 regions globally and keep 5 copies
           of everything for redundancy. We'll also use the most expensive
           database instances for maximum performance."
```

**Why It's Wrong:**
- Shows lack of business awareness
- Could make solution economically unfeasible
- Demonstrates poor engineering judgment
- Missing real-world constraints

**✅ Correct Approach:**
```
Good Candidate: "For cost optimization, I'll consider:
                 - Tiered storage (hot, warm, cold data)
                 - Regional deployment based on user distribution
                 - Auto-scaling to match demand
                 - Reserved instances for predictable workloads
                 - Cost monitoring and budgets"
```

**Key Lesson:** Always consider cost implications and optimize for business value.

## ⚙️ Technology and Scaling Mistakes

### ❌ Mistake #13: Wrong Technology Choices

**What It Looks Like:**
```
Candidate: "I'll use MongoDB for everything because it's NoSQL and scalable."
Interviewer: "What about ACID transactions for financial data?"
Candidate: "MongoDB can handle that too."
```

**Why It's Wrong:**
- One-size-fits-all approach
- Doesn't understand technology trade-offs
- May choose trendy over appropriate
- Shows lack of practical experience

**✅ Correct Approach:**
```
Good Candidate: "For technology choices:
                 - PostgreSQL for transactional data (ACID compliance)
                 - Redis for caching (fast in-memory access)
                 - Elasticsearch for search (full-text capabilities)
                 - S3 for file storage (durability and cost-effectiveness)
                 Each choice is based on specific requirements."
```

**Key Lesson:** Choose technologies based on specific requirements, not popularity.

### ❌ Mistake #14: Premature Optimization

**What It Looks Like:**
```
Candidate: "We need to optimize for microsecond latency, so I'll use
           custom C++ services, in-memory databases, and specialized
           hardware for this simple CRUD application."
```

**Why It's Wrong:**
- Adds unnecessary complexity
- Increases development and maintenance costs
- May not address actual bottlenecks
- Shows poor prioritization

**✅ Correct Approach:**
```
Good Candidate: "I'll start with standard technologies and optimize based on
                 actual performance measurements. If we hit bottlenecks,
                 I'll profile the system and optimize the specific areas
                 that need improvement."
```

**Key Lesson:** Optimize based on actual needs and measurements, not assumptions.

### ❌ Mistake #15: Ignoring Single Points of Failure

**What It Looks Like:**
```
Design: [Single database, single load balancer, single application server]
Interviewer: "What happens if your database fails?"
Candidate: "We'll restore from backup."
```

**Why It's Wrong:**
- Poor availability design
- Long recovery times
- Shows lack of reliability thinking
- Could violate SLA requirements

**✅ Correct Approach:**
```
Good Candidate: "To eliminate single points of failure:
                 - Database: Master-slave replication with automatic failover
                 - Load balancer: Multiple load balancers with health checks
                 - Application: Multiple instances across availability zones
                 - Network: Multiple network paths and providers"
```

**Key Lesson:** Identify and eliminate single points of failure for critical systems.

### ❌ Mistake #16: Poor Capacity Estimation

**What It Looks Like:**
```
Candidate: "We'll need about 100 servers."
Interviewer: "How did you calculate that?"
Candidate: "It seems like a good number."
```

**Why It's Wrong:**
- No mathematical basis
- Can't validate design decisions
- Shows lack of analytical thinking
- May lead to over/under-provisioning

**✅ Correct Approach:**
```
Good Candidate: "Let me calculate capacity:
                 - 1M users, 10% concurrent = 100K concurrent users
                 - Each user makes 10 requests/hour = 278 QPS
                 - Each server handles 1000 QPS = need 1 server
                 - With 3x peak factor = 3 servers
                 - With redundancy = 6 servers minimum"
```

**Key Lesson:** Always show your math and reasoning for capacity estimates.

## 🔍 Operational and Monitoring Mistakes

### ❌ Mistake #17: Forgetting About Monitoring and Observability

**What It Looks Like:**
```
Candidate: [Presents complete system design]
Interviewer: "How would you monitor this system?"
Candidate: "We could add some logging..."
```

**Why It's Wrong:**
- Monitoring is critical for production systems
- Shows lack of operational experience
- Makes debugging and optimization difficult
- Could lead to poor system reliability

**✅ Correct Approach:**
```
Good Candidate: "For monitoring and observability:
                 - Metrics: CPU, memory, disk, network, application metrics
                 - Logging: Structured logs with correlation IDs
                 - Tracing: Distributed tracing for request flows
                 - Alerting: SLA-based alerts with escalation
                 - Dashboards: Real-time system health visualization"
```

**Key Lesson:** Design monitoring and observability from the beginning.

### ❌ Mistake #18: Not Considering Data Migration and Deployment

**What It Looks Like:**
```
Interviewer: "How would you migrate from the old system to this new design?"
Candidate: "We'd just switch over during a maintenance window."
```

**Why It's Wrong:**
- Ignores migration complexity
- Could cause extended downtime
- Risk of data loss or corruption
- Shows lack of production experience

**✅ Correct Approach:**
```
Good Candidate: "For migration strategy:
                 - Dual-write to both old and new systems
                 - Gradual traffic migration with feature flags
                 - Data validation and consistency checks
                 - Rollback plan if issues arise
                 - Blue-green deployment for zero downtime"
```

**Key Lesson:** Always consider how to migrate and deploy safely.

### ❌ Mistake #19: Inadequate Error Handling

**What It Looks Like:**
```
Candidate: "If there's an error, we'll return a 500 status code."
Interviewer: "What about different types of errors?"
Candidate: "They're all just errors."
```

**Why It's Wrong:**
- Poor user experience
- Difficult to debug issues
- No differentiation between error types
- Missing graceful degradation

**✅ Correct Approach:**
```
Good Candidate: "For error handling:
                 - Specific error codes for different scenarios
                 - Graceful degradation when possible
                 - Retry logic with exponential backoff
                 - Circuit breakers for downstream failures
                 - Detailed error logging for debugging"
```

**Key Lesson:** Design comprehensive error handling and recovery strategies.

## 🎯 Interview-Specific Mistakes

### ❌ Mistake #20: Not Asking Questions

**What It Looks Like:**
```
Candidate: [Presents design without any questions]
Interviewer: "Do you have any questions for me?"
Candidate: "No, I think I covered everything."
```

**Why It's Wrong:**
- Misses opportunity to show curiosity
- Doesn't demonstrate collaborative thinking
- May miss important context
- Shows lack of engagement

**✅ Correct Approach:**
```
Good Candidate: "I have a few questions:
                 - What's the team's experience with the technologies I proposed?
                 - Are there any existing systems I should integrate with?
                 - What are the main technical challenges the team faces?
                 - How does this system fit into the broader architecture?"
```

**Key Lesson:** Always ask thoughtful questions to show engagement and curiosity.

## 🎯 Mistake Prevention Checklist

### ✅ Pre-Interview Preparation
```
□ Practice requirement clarification questions
□ Review capacity estimation formulas
□ Study common system design patterns
□ Understand technology trade-offs
□ Practice time management
□ Prepare thoughtful questions to ask
```

### ✅ During Interview Execution
```
□ Start with requirement clarification
□ Think out loud continuously
□ Start simple, add complexity gradually
□ Validate design against requirements
□ Consider failure scenarios
□ Discuss monitoring and operations
□ Ask for feedback and incorporate it
□ Manage time effectively
```

### ✅ Post-Interview Reflection
```
□ Identify what went well
□ Note areas for improvement
□ Review any questions you couldn't answer
□ Plan follow-up learning
□ Send thank you note
```

## 🎯 Summary: Common Mistakes to Avoid

### Critical Success Factors
- **Communication First** - Clarify requirements and think out loud
- **Start Simple** - Begin with high-level design, add complexity gradually
- **Consider Failures** - Design for resilience and operational reality
- **Show Your Math** - Provide concrete numbers and reasoning
- **Ask Questions** - Demonstrate curiosity and collaboration
- **Manage Time** - Cover all required topics within time limits

### Red Flags That Signal Problems
- Jumping to solutions without understanding requirements
- Over-engineering for unrealistic scale
- Ignoring failure scenarios and edge cases
- Poor time management and rushed explanations
- Not incorporating interviewer feedback
- Missing security and operational considerations

### Keys to Interview Success
- **Systematic Approach** - Follow consistent methodology
- **Clear Communication** - Explain reasoning and trade-offs
- **Technical Depth** - Show appropriate level of detail
- **Real-world Awareness** - Consider practical constraints
- **Collaborative Spirit** - Work with interviewer, not against them

**By avoiding these common mistakes, you significantly increase your chances of interview success!** 🚀

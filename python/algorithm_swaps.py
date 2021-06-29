# Function to count Inversions
def howManySwaps(arr, n):
    # A temp_arr is created to store
    # sorted array in merge function
    temp_arr = [0]*n
    return _merge_sort(arr, temp_arr, 0, n-1)


def _merge_sort(array, temp_array, l, r):
    inv_count = 0
    # We will make a recursive call if and only if
    # we have more than one elements
    if l < r:
        # mid is calculated to divide the array into two subarrays
        # Floor division is must in case of python
        mid = (l + r)//2
        # recursively calculate inversions on left subarray
        inv_count += _merge_sort(array, temp_array, l, mid)
        # recursively calculate inversions on right subarray
        inv_count += _merge_sort(array, temp_array, mid + 1, r)
        # merge the 2 sub-arrays
        inv_count += merge(array, temp_array, l, mid, r)
    return inv_count


# function to merge 2 subarrays
def merge(arr, temp_arr, left, mid, right):
    i = left         # Starting index of left subarray
    j = mid + 1      # Starting index of right subarray
    k = left         # Starting index of to be sorted subarray
    inv_count = 0

    while i <= mid and j <= right:
        # There will be no inversion if arr[i] <= arr[j]
        if arr[i] <= arr[j]:
            temp_arr[k] = arr[i]
            k += 1
            i += 1
        else:
            # Inversion will occur.
            temp_arr[k] = arr[j]
            inv_count += (mid-i + 1)
            k += 1
            j += 1
    while i <= mid:
        temp_arr[k] = arr[i]
        k += 1
        i += 1
    while j <= right:
        temp_arr[k] = arr[j]
        k += 1
        j += 1
    for loop_var in range(left, right + 1):
        arr[loop_var] = temp_arr[loop_var]

    return inv_count


n = int(input())
arr = []
for i in range(0, n):
    arr.append(int(input()))
print(howManySwaps(arr, n))

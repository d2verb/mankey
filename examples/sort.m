selectionSort = fn(array) {
  i = 0
  while (i < len(array)) {
    j = i + 1
    while (j < len(array)) {
      if (array[j] < array[i]) {
        tmp = array[j]
        array[j] = array[i]
        array[i] = tmp
      }
      j = j + 1
    }
    i = i + 1
  }
}

main = fn() {
  array = [5, 3, 9, 1, 11, 13, 2, 52, 3, 100]
  selectionSort(array)

  puts(array)
}

main()

map = fn(arr, f) {
  new_arr = copy(arr)
  i = 0
  while (i < len(arr)) {
    new_arr[i] = f(new_arr[i])
    i = i + 1
  }
  return new_arr
}

main = fn() {
  arr = [1, 2, 3, 4, 5]
  new_arr = map(arr, fn(x) { x * x * x })
  println(arr)
  println(new_arr)
}

main()

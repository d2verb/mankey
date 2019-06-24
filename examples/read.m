main = fn() {
  fd = open("./read.m", "r")
  nb = 1
  while (nb > 0) {
    res = read(fd, 10)
    print(res.content)
    nb = res.nb
  }
  close(fd)
}

main()

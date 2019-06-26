readln = fn() {
  got = {"content": ""}
  line = ""
  while (got.content != "\n") {
    got = read(0, 1)
    line = line + got.content
  }
  return line
}

print(readln())

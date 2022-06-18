EXE=market

run:
	go build -o $(EXE)
	./$(EXE)

install:
	go build -o $(EXE)
	mv $(EXE) /bin/$(EXE)
platform:
	go build -o bin/$(EXE)
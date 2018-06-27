SRC_DIR := ./../../src

CC := gcc
CFLAGS += -Wall -Os -s -I. -I$(SRC_DIR)

.PHONY: all
all:
	go install github.com/jlubawy/go-ctlog/...
	cmodule json . > cmodule_indices.json
	cmodule header cmodule_indices.json > cmodule_indices.h
	ctlog json cmodule_indices.json > ctlog_dict.json
	$(CC) $(CFLAGS) -o main main.c
	./main > main.txt
	$(CC) $(CFLAGS) -o main_ctlog $(SRC_DIR)/ctlog.c main_ctlog.c
	./main_ctlog | ctlog log ctlog_dict.json > main_ctlog.txt

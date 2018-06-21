################################################################################
# Makefile for Arduino Uno projects
################################################################################

# Project name
APP_TARGET := ctlog_example_arduino

# ComPort to use for installing the binary, change as needed
AVR_DUDE_COM_PORT ?= COM3

# Use CMD as the default shell if the OS is Windows
ifeq ($(OS),Windows_NT)
SHELL = cmd
endif

# MCU and definitions for the Arduino Uno
MMCU ?= atmega328p
DEFINES := -DF_CPU=16000000L
DEFINES += -DARDUINO=10603
DEFINES += -DARDUINO_AVR_UNO
DEFINES += -DARDUINO_ARCH_AVR

# Arduino installation directories, change as needed
DIR_ARDUINO := C:/PROGRA~2/Arduino
DIR_AVR_BIN := $(DIR_ARDUINO)/hardware/tools/avr/bin
DIR_AVR_INCLUDE := $(DIR_ARDUINO)/hardware/tools/avr/avr/include
DIR_AVR_ETC := $(DIR_ARDUINO)/hardware/tools/avr/etc

DIR_ROOT := $(abspath $(CURDIR))
DIR_BIN := $(DIR_ROOT)/bin
DIR_OBJ := $(DIR_ROOT)/obj
DIR_SRC := $(DIR_ROOT)/src
DIR_SRC_GEN := $(DIR_ROOT)/src_gen
DIR_SRC_CTLOG := $(abspath $(DIR_ROOT)/../../src)

AR := $(DIR_AVR_BIN)/avr-ar
CC := $(DIR_AVR_BIN)/avr-gcc
CXX := $(DIR_AVR_BIN)/avr-g++

# Include directories, change as needed
INCLUDES := $(addprefix -I,.)
INCLUDES += $(addprefix -I,$(DIR_SRC))
INCLUDES += $(addprefix -I,$(DIR_SRC_GEN))
INCLUDES += $(addprefix -I,$(DIR_SRC_CTLOG))
INCLUDES += $(addprefix -I,$(DIR_AVR_INCLUDE))

CFLAGS = -g -Os -Wl,-Map,$@.map -Wl,--gc-sections -Wl,-u,vfprintf -lprintf_flt -lm -fno-exceptions -ffunction-sections -fdata-sections -MMD -mmcu=$(MMCU) $(DEFINES) $(INCLUDES)
CXXFLAGS = -g -Os -Wl,-Map,$@.map -Wl,--gc-sections -Wl,-u,vfprintf -lprintf_flt -lm -fno-exceptions -ffunction-sections -fdata-sections -MMD -mmcu=$(MMCU) $(DEFINES) $(INCLUDES)

# C and C++ source files, change as needed
SRC_C := $(wildcard $(DIR_SRC)/*.c)
SRC_C += $(wildcard $(DIR_SRC_CTLOG)/*.c)
SRC_CXX =

# Object files
OBJ_C := $(addprefix $(DIR_OBJ)/,$(notdir $(SRC_C:.c=.c.o)))
OBJ_CXX := $(addprefix $(DIR_OBJ)/,$(notdir $(SRC_CXX:.c=.c.o)))

################################################################################
# Application Recipes
################################################################################
OUT_ELF = $(addprefix $(DIR_BIN)/,$(addsuffix .elf,$(APP_TARGET)))
OUT_EEP = $(addprefix $(DIR_BIN)/,$(addsuffix .eep,$(APP_TARGET)))
OUT_HEX = $(addprefix $(DIR_BIN)/,$(addsuffix .hex,$(APP_TARGET)))

################################################################################
# Targets
################################################################################
.PHONY: all
all:
	go install github.com/jlubawy/go-ctlog/...
	$(MAKE) dirs
	$(MAKE) generate_cmodule
	$(MAKE) generate_ctlog
	$(MAKE) $(OUT_HEX)

.PHONY: install
install:
	$(DIR_AVR_BIN)/avrdude -C$(DIR_AVR_ETC)/avrdude.conf -v -p$(MMCU) -carduino -P$(AVR_DUDE_COM_PORT) -b115200 -D -Uflash:w:$(OUT_HEX):i

.PHONY: generate_cmodule
generate_cmodule:
	cmodule json ./src > $(DIR_SRC_GEN)/cmodule_indices.json
	cmodule header $(DIR_SRC_GEN)/cmodule_indices.json > $(DIR_SRC_GEN)/cmodule_indices.h

.PHONY: generate_ctlog
generate_ctlog:
	ctlog json $(DIR_SRC_GEN)/cmodule_indices.json > $(DIR_BIN)/ctlog_dict.json

.PHONY: clean
clean:
ifeq ($(SHELL),cmd)
	if exist "bin" for /d %%G in ("bin") do rd /s /q "%%G"
	if exist "obj" for /d %%G in ("obj") do rd /s /q "%%G"
	if exist "src_gen" for /d %%G in ("src_gen") do rd /s /q "%%G"
else
	rm -rf bin
	rm -rf obj
	rm -rf src_gen
endif

.PHONY: dirs
dirs:
ifeq ($(SHELL),cmd)
	if not exist "bin" mkdir bin
	if not exist "obj" mkdir obj
	if not exist "src_gen" mkdir src_gen
else
	mkdir -p bin
	mkdir -p obj
	mkdir -p src_gen
endif

# Object Files
$(DIR_OBJ)/%.c.o: $(DIR_SRC)/%.c
	$(info $@)
	$(CC) $(CFLAGS) -o $@ -c $<
$(DIR_OBJ)/%.c.o: $(DIR_SRC_CTLOG)/%.c
	$(info $@)
	$(CC) $(CFLAGS) -o $@ -c $<

$(DIR_OBJ)/%.cpp.o: $(DIR_SRC)/%.cpp
	$(info $@)
	$(CXX) $(CXXFLAGS) -fno-threadsafe-statics -o $@ -c $<
$(DIR_OBJ)/%.cpp.o: $(DIR_SRC_CTLOG)/%.cpp
	$(info $@)
	$(CXX) $(CXXFLAGS) -fno-threadsafe-statics -o $@ -c $<

# Output Files
$(OUT_ELF): $(OBJ_C) $(OBJ_CXX)
	$(info $@)
	$(CC) $(CFLAGS) -o $@ $^
	$(DIR_AVR_BIN)/avr-size --mcu=$(MMCU) -C $@
	$(DIR_AVR_BIN)/avr-size --mcu=$(MMCU) -C $@ > $(DIR_BIN)/memory_usage.txt
	$(DIR_AVR_BIN)/avr-objdump --disassemble --syms $@ > $@.lst

$(OUT_EEP) $(OUT_HEX): $(OUT_ELF)
	$(info $@)
	$(DIR_AVR_BIN)/avr-objcopy -O ihex -j .eeprom --set-section-flags=.eeprom=alloc,load --no-change-warnings --change-section-lma .eeprom=0 $< $@
	$(DIR_AVR_BIN)/avr-objcopy -O ihex -R .eeprom $< $@

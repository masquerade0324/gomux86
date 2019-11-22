TARGET := test
OBJS   := crt0.o test.o

CC := gcc
AS := nasm
LD := ld

CFLAGS  += -m32 -nostdlib -fno-asynchronous-unwind-tables -fno-pie -g -fno-stack-protector
LDFLAGS += -m elf_i386 --entry=start --oformat=binary -Ttext 0x7c00

.PHONY: all
all:
	make $(TARGET)

%.o: %.c Makefile
	$(CC) $(CFLAGS) -c $<

%.o: %.asm Makefile
	$(AS) -f elf $<

$(TARGET): $(OBJS) Makefile
	$(LD) $(LDFLAGS) -o $@ $(OBJS)

clean:
	$(RM) $(TARGET) $(OBJS)
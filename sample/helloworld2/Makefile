TARGET := helloworld2

AS := nasm

.PHONY: all
all:
	make $(TARGET)

%.o: %.asm Makefile
	$(AS) -f elf $<

$(TARGET): $(TARGET).asm
	$(AS) -f bin -o $(TARGET) $<

clean:
	$(RM) $(TARGET)

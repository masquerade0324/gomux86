target = in out select

AS = nasm

.PHONY: all
all:
	make $(target)

%: %.asm Makefile
	$(AS) -f bin -o $@ $<

.PHONY: clean
clean:
	rm -f $(target)

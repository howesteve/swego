# Compile first the Swiss Ephemeris as static library!
# see ../../swisseph/Makefile

CC = cc
CFLAGS += -std=c99 -g -Wall -fPIC
LDFLAGS += -Wl -dead_strip

.o:
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ $? -L../../swisseph -lswe -lm
.c.o:
	$(CC) $(CFLAGS) -c -I../../swisseph $<

swerker-stdio: swerker.o msgpuck.o handlers.o tr-stdio.o swex.o
	$(CC) $(CFLAGS) $(LDFLAGS) -L../../swisseph -lswe -lm \
		-o swerker-stdio swerker.o msgpuck.o handlers.o tr-stdio.o swex.o

clean:
	rm -f *.o swerker-stdio

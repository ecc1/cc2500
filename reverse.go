package cc2500

/*
From: cooper@decwrl.ARPA (Eric Cooper)
Message-Id: <8504242012.AA09926@lewis.ARPA>
Date: 24 Apr 1985 1211-PST (Wednesday)
To: cooper@BERKELEY
Subject: fast bit reversal

The following O(n) algorithm for bit-reversing n numbers
(assuming that a (log n)-bit XOR takes unit time)
was suggested by Andrei Broder.
By comparison, shifting and masking algorithms are O(n log n).

First, observe that the bit pattern for i XOR i+1 always consists of
some number of 0s on the left followed by some number of 1s on the right.
The 1s represent how far the carry ripples in going from i to i+1.

Let A[0..2^n-1] be the table we are trying to compute.
The first step is to seed it with the reversals of
each of the special patterns above.

Now observe that we can simulate the "reverse carry" in going
from A[i] to A[i+1] using the bit pattern A[i XOR i+1].
Thus,
      A[i+1] := A[i] XOR A[i XOR i+1]
where the first term on the right comes from the previous iteration
through a loop, and the second term comes from the seeding operation.
------------------------------------------------------------------------

MODULE BitReversal;

FROM BitOperations IMPORT BitXor;

VAR
  Table: ARRAY [0..255] OF [0..255];

PROCEDURE ReverseBits;
  VAR
    i, j, x: INTEGER;
  BEGIN
    (* seed table *)
    i := 0; x := 0;
    j := 128;
    LOOP
	Table[i] := x;
	IF j = 0 THEN EXIT END;
	i := 2*i+1; x := x+j; j := j DIV 2;
    END;
    (* fill in rest of table *)
    FOR i := 1 TO 253 DO
      Table[i+1] := BitXor(Table[i],Table[BitXor(i,i+1)]);
    END;
  END ReverseBits;

END BitReversal.
*/

var reverseBits [256]byte

func init() {
	// Seed the table.
	i := 0
	x := byte(0)
	j := byte(128)
	for {
		reverseBits[i] = x
		if j == 0 {
			break
		}
		i = 2*i + 1
		x += j
		j /= 2
	}
	// Fill in the rest of the table.
	for i := 1; i < 254; i++ {
		reverseBits[i+1] = reverseBits[i] ^ reverseBits[i^(i+1)]
	}
}

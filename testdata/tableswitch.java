/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
class tableswitch {
	// simple test class that will use the TABLESWITCH bytecode
	public static void main(String[] args) {
		int i;
    	switch (args.length) {
        	case 0:  i =  0; break;
        	case 1:  i =  1; break;
        	case 2:  i =  2; break;
        	default: i = -1; break;
    	}
		System.out.printf("Value based on args is: %d\n", i );
	}
}
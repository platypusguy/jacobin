class lookupswitch {
    // class to test functioning of lookupswitch bytecode
    public static void main(String[] args) {
        var len = args.length;
        switch (len) {
            case 0: System.out.println("zero args"); break;
            case -100: System.out.println("100 args"); break;
            case 250: System.out.println("250 args"); break;
            default: System.out.println("args != -100, 0, or 250"); break;
        }
    }
}
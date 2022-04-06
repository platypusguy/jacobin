import java.util.List;
import java.util.ArrayList;

class ListTest {

	public static void main(String[] args) {
		List<String> greetings = new ArrayList<>();
		greetings.add("Hello!");
		greetings.add("Hi!");
		greetings.add("Welcome");

		for (String greeting : greetings) {
			System.out.println(greeting);
		}
	}

}

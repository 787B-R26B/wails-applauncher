from calculator import add

try:
    num1_str = input("Enter the first number: ")
    num1 = float(num1_str)
    num2_str = input("Enter the second number: ")
    num2 = float(num2_str)
    result = add(num1, num2)
    print(f"The sum is: {result}")
except ValueError:
    print("Invalid input. Please enter numbers only.")
except EOFError:
    print("No input received. Exiting.")

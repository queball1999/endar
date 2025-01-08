from argon2 import PasswordHasher

# You may need to run pip install argon2-cffi
# This is not included in requirements.txt for security reasons

def hash_password(password: str) -> str:
    # Initialize the PasswordHasher
    ph = PasswordHasher()
    
    # Hash the password
    hashed_password = ph.hash(password)
    return hashed_password

if __name__ == "__main__":
    # Prompt the user for the password
    password = input("Enter your password: ")
    
    # Prompt to see if this is for Docker Compose
    use_docker_compose = None
    while use_docker_compose not in ["yes", "no"]:
        use_docker_compose = input("Will you be using this hash in a Docker Compose file? (yes/no): ").strip().lower()
    
    # Hash the password
    hashed_password = hash_password(password)
    
    # Replace single $ with $$ if using Docker Compose
    if use_docker_compose == "yes":
        hashed_password = hashed_password.replace("$", "$$")
    
    # Output the hashed password
    print("Hashed Password: {}".format(hashed_password))

import logging
import os

# Ensure the logs directory exists
log_dir = os.path.abspath(os.path.dirname(__file__))
log_file = os.path.join(log_dir, "app.log")

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,  # Change to logging.INFO in production
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[
        logging.FileHandler(log_file),  # Write logs to a file
        logging.StreamHandler()  # Also print logs to stdout
    ]
)

logger = logging.getLogger(__name__)

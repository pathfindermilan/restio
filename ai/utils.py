import base64

def image_to_base64(image_data: bytes) -> str:
    """
    Converts binary image data to a base64-encoded string.

    Args:
        image_data (bytes): The binary data of the image.

    Returns:
        str: The base64-encoded string of the image.
    """
    try:
        base64_string = base64.b64encode(image_data).decode("utf-8")
        return base64_string
    except Exception as e:
        return f"An error occurred: {str(e)}"

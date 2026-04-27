from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC

from .base_page import BasePage


class CartPage(BasePage):
    heading = (By.TAG_NAME, "h2")
    empty_message = (By.XPATH, "//p[normalize-space()='Koszyk jest pusty']")
    cart_items = (By.CSS_SELECTOR, "section ul li")
    remove_buttons = (By.XPATH, "//button[normalize-space()='Usuń']")

    def wait_for_loaded(self) -> None:
        self.wait.until(EC.text_to_be_present_in_element(self.heading, "Koszyk"))

    def is_empty_visible(self) -> bool:
        return len(self.driver.find_elements(*self.empty_message)) > 0

    def get_item_texts(self) -> list[str]:
        return [item.text for item in self.driver.find_elements(*self.cart_items)]

    def remove_first_item(self) -> None:
        self.wait.until(EC.element_to_be_clickable(self.remove_buttons)).click()

    def wait_for_empty(self) -> None:
        self.wait.until(EC.visibility_of_element_located(self.empty_message))

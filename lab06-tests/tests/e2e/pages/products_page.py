from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC

from .base_page import BasePage


class ProductsPage(BasePage):
    heading = (By.TAG_NAME, "h2")
    items = (By.CSS_SELECTOR, "section ul li")
    add_buttons = (By.XPATH, "//button[normalize-space()='Dodaj do koszyka']")

    def wait_for_loaded(self) -> None:
        self.wait.until(EC.text_to_be_present_in_element(self.heading, "Produkty"))
        self.wait.until(EC.presence_of_all_elements_located(self.items))

    def get_item_texts(self) -> list[str]:
        return [item.text for item in self.wait.until(EC.presence_of_all_elements_located(self.items))]

    def add_first_product(self) -> str:
        first_item_text = self.get_item_texts()[0]
        self.wait.until(EC.element_to_be_clickable(self.add_buttons)).click()
        return first_item_text.split(" - ")[0]

    def get_add_button_count(self) -> int:
        return len(self.wait.until(EC.presence_of_all_elements_located(self.add_buttons)))

from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC

from .base_page import BasePage


class AppPage(BasePage):
    products_link = (By.LINK_TEXT, "Produkty")
    cart_link = (By.PARTIAL_LINK_TEXT, "Koszyk")
    payments_link = (By.LINK_TEXT, "Płatności")
    app_title = (By.TAG_NAME, "h1")

    def wait_for_app_shell(self) -> None:
        self.wait.until(EC.text_to_be_present_in_element(self.app_title, "Sklep"))

    def go_to_products(self) -> None:
        self.wait.until(EC.element_to_be_clickable(self.products_link)).click()

    def go_to_cart(self) -> None:
        self.wait.until(EC.element_to_be_clickable(self.cart_link)).click()

    def go_to_payments(self) -> None:
        self.wait.until(EC.element_to_be_clickable(self.payments_link)).click()

    def get_cart_badge_text(self) -> str:
        return self.wait.until(EC.visibility_of_element_located(self.cart_link)).text

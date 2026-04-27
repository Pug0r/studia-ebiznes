from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC

from .base_page import BasePage


class PaymentsPage(BasePage):
    heading = (By.TAG_NAME, "h2")
    email_input = (By.CSS_SELECTOR, "input[type='email']")
    amount_input = (By.CSS_SELECTOR, "input[type='number']")
    submit_button = (By.XPATH, "//button[normalize-space()='Wyślij']")
    success_message = (By.XPATH, "//p[normalize-space()='Płatność zapisana']")
    saved_payments_items = (By.CSS_SELECTOR, "section ul li")

    def wait_for_loaded(self) -> None:
        self.wait.until(EC.text_to_be_present_in_element(self.heading, "Płatności"))

    def get_saved_payment_count(self) -> int:
        return len(self.driver.find_elements(*self.saved_payments_items))

    def send_payment(self, email: str, amount: str) -> None:
        email_field = self.wait.until(EC.visibility_of_element_located(self.email_input))
        amount_field = self.wait.until(EC.visibility_of_element_located(self.amount_input))
        email_field.clear()
        email_field.send_keys(email)
        amount_field.clear()
        amount_field.send_keys(amount)
        self.wait.until(EC.element_to_be_clickable(self.submit_button)).click()

    def wait_for_success(self) -> None:
        self.wait.until(EC.visibility_of_element_located(self.success_message))

    def has_saved_payment(self, email: str) -> bool:
        items = self.driver.find_elements(*self.saved_payments_items)
        return any(email in item.text for item in items)

    def is_form_valid(self) -> bool:
        form = self.driver.find_element(By.TAG_NAME, "form")
        return self.driver.execute_script("return arguments[0].checkValidity();", form)

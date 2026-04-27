from uuid import uuid4

from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait

from tests.e2e.pages.app_page import AppPage
from tests.e2e.pages.cart_page import CartPage
from tests.e2e.pages.payments_page import PaymentsPage
from tests.e2e.pages.products_page import ProductsPage


def test_app_title_visible(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()

    assert driver.find_element(By.TAG_NAME, "h1").text == "Sklep"


def test_redirects_to_products_view_on_open(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()

    WebDriverWait(driver, wait_timeout).until(EC.text_to_be_present_in_element((By.TAG_NAME, "h2"), "Produkty"))
    assert "/products" in driver.current_url


def test_navigation_contains_products_link(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()

    assert driver.find_element(By.LINK_TEXT, "Produkty").is_displayed()


def test_navigation_contains_cart_link(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()

    assert driver.find_element(By.PARTIAL_LINK_TEXT, "Koszyk").is_displayed()


def test_navigation_contains_payments_link(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()

    assert driver.find_element(By.LINK_TEXT, "Płatności").is_displayed()


def test_products_list_is_visible(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    products_page = ProductsPage(driver, wait_timeout)
    products_page.wait_for_loaded()

    assert len(products_page.get_item_texts()) > 0


def test_products_contains_laptop_item(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    products_page = ProductsPage(driver, wait_timeout)
    products_page.wait_for_loaded()

    assert any("Laptop" in text for text in products_page.get_item_texts())


def test_each_product_has_add_button(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    products_page = ProductsPage(driver, wait_timeout)
    products_page.wait_for_loaded()

    assert products_page.get_add_button_count() == len(products_page.get_item_texts())


def test_cart_badge_starts_at_zero(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()

    assert "Koszyk (0)" in app_page.get_cart_badge_text()


def test_cart_view_shows_empty_message_initially(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    app_page.go_to_cart()
    cart_page = CartPage(driver, wait_timeout)
    cart_page.wait_for_loaded()

    assert cart_page.is_empty_visible()


def test_adding_product_updates_cart_badge(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    products_page = ProductsPage(driver, wait_timeout)
    products_page.wait_for_loaded()
    products_page.add_first_product()

    assert "Koszyk (1)" in app_page.get_cart_badge_text()


def test_added_product_is_visible_in_cart(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    products_page = ProductsPage(driver, wait_timeout)
    products_page.wait_for_loaded()
    product_name = products_page.add_first_product()

    app_page.go_to_cart()
    cart_page = CartPage(driver, wait_timeout)
    cart_page.wait_for_loaded()
    cart_items = cart_page.get_item_texts()

    assert any(product_name in item_text for item_text in cart_items)


def test_removing_product_empties_cart(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    products_page = ProductsPage(driver, wait_timeout)
    products_page.wait_for_loaded()
    products_page.add_first_product()

    app_page.go_to_cart()
    cart_page = CartPage(driver, wait_timeout)
    cart_page.wait_for_loaded()
    cart_page.remove_first_item()
    cart_page.wait_for_empty()

    assert cart_page.is_empty_visible()


def test_payments_view_is_visible(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    app_page.go_to_payments()
    payments_page = PaymentsPage(driver, wait_timeout)
    payments_page.wait_for_loaded()

    assert driver.find_element(By.TAG_NAME, "h2").text == "Płatności"


def test_payments_form_has_required_email_field(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    app_page.go_to_payments()

    email_input = WebDriverWait(driver, wait_timeout).until(EC.visibility_of_element_located((By.CSS_SELECTOR, "input[type='email']")))
    assert email_input.get_attribute("required") is not None


def test_payments_form_has_numeric_amount_field(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    app_page.go_to_payments()

    amount_input = WebDriverWait(driver, wait_timeout).until(EC.visibility_of_element_located((By.CSS_SELECTOR, "input[type='number']")))
    assert amount_input.get_attribute("type") == "number"


def test_submitting_valid_payment_shows_success_message(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    app_page.go_to_payments()
    payments_page = PaymentsPage(driver, wait_timeout)
    payments_page.wait_for_loaded()

    unique_email = f"selenium-{uuid4().hex[:8]}@test.local"
    payments_page.send_payment(unique_email, "123.45")
    payments_page.wait_for_success()

    assert driver.find_element(By.XPATH, "//p[normalize-space()='Płatność zapisana']").is_displayed()


def test_submitting_valid_payment_adds_item_to_saved_payments(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    app_page.go_to_payments()
    payments_page = PaymentsPage(driver, wait_timeout)
    payments_page.wait_for_loaded()

    before_count = payments_page.get_saved_payment_count()
    unique_email = f"selenium-{uuid4().hex[:8]}@test.local"
    payments_page.send_payment(unique_email, "99.99")
    payments_page.wait_for_success()

    WebDriverWait(driver, wait_timeout).until(lambda current_driver: payments_page.get_saved_payment_count() == before_count + 1)
    assert payments_page.has_saved_payment(unique_email)


def test_payment_form_is_invalid_for_wrong_email(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    app_page.go_to_payments()
    payments_page = PaymentsPage(driver, wait_timeout)
    payments_page.wait_for_loaded()

    payments_page.send_payment("wrong-email", "10")
    assert payments_page.is_form_valid() is False


def test_amount_input_has_min_zero(driver, wait_timeout):
    app_page = AppPage(driver, wait_timeout)
    app_page.wait_for_app_shell()
    app_page.go_to_payments()

    amount_input = WebDriverWait(driver, wait_timeout).until(EC.visibility_of_element_located((By.CSS_SELECTOR, "input[type='number']")))
    assert amount_input.get_attribute("min") == "0"

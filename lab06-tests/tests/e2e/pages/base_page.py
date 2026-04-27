from selenium.webdriver.remote.webdriver import WebDriver
from selenium.webdriver.support.ui import WebDriverWait


class BasePage:
    def __init__(self, driver: WebDriver, wait_timeout: int):
        self.driver = driver
        self.wait = WebDriverWait(driver, wait_timeout)

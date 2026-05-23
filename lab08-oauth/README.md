Należy skonfigurować klienta Oauth2 (4.0). Dane o użytkowniku wraz z
tokenem powinny być przechowywane po stronie bazy serwera, a nowy
token (inny niż ten od dostawcy) powinien zostać wysłany do klienta
(React). Można zastosować mechanizm sesji lub inny dowolny (5.0).
Zabronione jest tworzenie klientów bezpośrednio po stronie React'a
wyłączając z komunikacji aplikację serwerową.

Prawidłowa komunikacja: react-sewer-dostawca-serwer(via return
uri)-react.

- [x] 3.0 logowanie przez aplikację serwerową (bez Oauth2)
- [x] 3.5 rejestracja przez aplikację serwerową (bez Oauth2)
- [x] 4.0 logowanie via Google OAuth2
- [ ] 4.5 logowanie via Facebook lub Github OAuth2
- [ ] 5.0 zapisywanie danych logowania OAuth2 po stronie serwera

Klucz należy uzyskać na:
- https://console.cloud.google.com/apis/dashboard,
- https://developers.facebook.com/

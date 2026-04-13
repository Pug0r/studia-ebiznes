import { Link, Navigate, Route, Routes } from 'react-router-dom'
import Products from './components/Products'
import Payments from './components/Payments'
import Cart from './components/Cart'
import { CartProvider, useCart } from './hooks/useCart'
import { PaymentsProvider, usePayments } from './hooks/usePayments'

export default function App() {
  return (
    <PaymentsProvider>
      <CartProvider>
        <AppContent />
      </CartProvider>
    </PaymentsProvider>
  )
}

function AppContent() {
  const { cartItems } = useCart()
  const { payments } = usePayments()

  return (
    <main>
      <h1>Sklep</h1>
      <nav>
        <Link to="/products">Produkty</Link>
        <Link to="/cart">Koszyk ({cartItems.length})</Link>
        <Link to="/payments">Płatności</Link>
      </nav>
      <Routes>
        <Route path="/products" element={<Products />} />
        <Route path="/cart" element={<Cart />} />
        <Route
          path="/payments"
          element={
            <div className="layout">
              <Payments />
              <section>
                <h2>Zapisane płatności</h2>
                <ul>
                  {payments.map((payment, index) => (
                    <li key={`${payment.email}-${payment.amount}-${index}`}>
                      {payment.email} - {payment.amount} PLN
                    </li>
                  ))}
                </ul>
              </section>
            </div>
          }
        />
        <Route path="*" element={<Navigate to="/products" replace />} />
      </Routes>
    </main>
  )
}
import { useEffect, useState } from 'react'
import { Link, Navigate, Route, Routes } from 'react-router-dom'
import Products from './components/Products'
import Payments from './components/Payments'
import Cart from './components/Cart'
import type { Product } from './components/Products'

const API_URL = import.meta.env.VITE_API_URL || '/api'

export type PaymentPayload = {
  email: string
  amount: number
}

export default function App() {
  const [cartItems, setCartItems] = useState<Product[]>([])
  const [message, setMessage] = useState('')
  const [payments, setPayments] = useState<PaymentPayload[]>([])

  async function loadPayments(): Promise<void> {
    const response = await fetch(`${API_URL}/payments`)
    if (!response.ok) {
      throw new Error('Błąd pobierania płatności')
    }
    const data: PaymentPayload[] = await response.json()
    setPayments(data)
  }

  useEffect(() => {
    loadPayments().catch(() => {
      setMessage('Nie udało się pobrać płatności')
    })
  }, [])

  function addToCart(product: Product): void {
    setCartItems((previous) => [...previous, product])
  }

  function removeFromCart(index: number): void {
    setCartItems((previous) => previous.filter((_, itemIndex) => itemIndex !== index))
  }

  async function handlePayment(data: PaymentPayload): Promise<void> {
    setMessage('Wysyłanie...')
    try {
      const response = await fetch(`${API_URL}/payments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      })
      if (!response.ok) {
        throw new Error('Błąd płatności')
      }
      await loadPayments()
      setMessage('Płatność zapisana')
    } catch {
      setMessage('Nie udało się wysłać płatności')
    }
  }

  return (
    <main>
      <h1>Sklep</h1>
      <nav>
        <Link to="/products">Produkty</Link>
        <Link to="/cart">Koszyk ({cartItems.length})</Link>
        <Link to="/payments">Płatności</Link>
      </nav>
      <Routes>
        <Route path="/products" element={<Products apiUrl={API_URL} onAddToCart={addToCart} />} />
        <Route path="/cart" element={<Cart items={cartItems} onRemove={removeFromCart} />} />
        <Route
          path="/payments"
          element={
            <div className="layout">
              <Payments onSubmit={handlePayment} message={message} />
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
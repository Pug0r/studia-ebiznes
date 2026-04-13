import { useEffect, useState } from 'react'
import Products from './components/Products'
import Payments from './components/Payments'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export type PaymentPayload = {
  email: string
  amount: number
}

export default function App() {
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
      <div className="layout">
        <div>
          <Products apiUrl={API_URL} />
          <Payments onSubmit={handlePayment} message={message} />
        </div>
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
    </main>
  )
}
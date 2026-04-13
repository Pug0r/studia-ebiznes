import { createContext, useContext, useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { API_URL } from '../config'
import type { PaymentPayload } from '../types'

type PaymentsContextValue = {
  payments: PaymentPayload[]
  message: string
  handlePayment: (data: PaymentPayload) => Promise<void>
}

const PaymentsContext = createContext<PaymentsContextValue | undefined>(undefined)

export function PaymentsProvider({ children }: { children: ReactNode }) {
  const [payments, setPayments] = useState<PaymentPayload[]>([])
  const [message, setMessage] = useState('')

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
    <PaymentsContext.Provider value={{ payments, message, handlePayment }}>
      {children}
    </PaymentsContext.Provider>
  )
}

export function usePayments(): PaymentsContextValue {
  const context = useContext(PaymentsContext)
  if (!context) {
    throw new Error('usePayments must be used within PaymentsProvider')
  }
  return context
}
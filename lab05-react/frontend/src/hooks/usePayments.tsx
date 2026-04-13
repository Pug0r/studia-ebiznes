import { createContext, useContext, useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import axios from 'axios'
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
    const response = await axios.get<PaymentPayload[]>(`${API_URL}/payments`)
    setPayments(response.data)
  }

  useEffect(() => {
    loadPayments().catch(() => {
      setMessage('Nie udało się pobrać płatności')
    })
  }, [])

  async function handlePayment(data: PaymentPayload): Promise<void> {
    setMessage('Wysyłanie...')
    try {
      await axios.post(`${API_URL}/payments`, data)
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
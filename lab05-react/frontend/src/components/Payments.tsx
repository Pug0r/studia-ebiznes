import { useState } from 'react'
import type { FormEvent } from 'react'
import type { PaymentPayload } from '../App'

type PaymentsProps = {
  onSubmit: (payload: PaymentPayload) => Promise<void>
  message: string
}

export default function Payments({ onSubmit, message }: PaymentsProps) {
  const [email, setEmail] = useState('')
  const [amount, setAmount] = useState('')

  async function handleSubmit(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault()
    await onSubmit({
      email,
      amount: Number(amount)
    })
  }

  return (
    <section>
      <h2>Płatności</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          required
        />
        <input
          type="number"
          placeholder="Kwota"
          value={amount}
          onChange={(event) => setAmount(event.target.value)}
          min="0"
          step="0.01"
          required
        />
        <button type="submit">Wyślij</button>
      </form>
      {message && <p>{message}</p>}
    </section>
  )
}
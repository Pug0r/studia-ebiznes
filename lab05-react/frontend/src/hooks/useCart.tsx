import { createContext, useContext, useState } from 'react'
import type { ReactNode } from 'react'
import type { Product } from '../types'

type CartContextValue = {
  cartItems: Product[]
  addToCart: (product: Product) => void
  removeFromCart: (index: number) => void
}

const CartContext = createContext<CartContextValue | undefined>(undefined)

export function CartProvider({ children }: { children: ReactNode }) {
  const [cartItems, setCartItems] = useState<Product[]>([])

  function addToCart(product: Product): void {
    setCartItems((previous) => [...previous, product])
  }

  function removeFromCart(index: number): void {
    setCartItems((previous) => previous.filter((_, itemIndex) => itemIndex !== index))
  }

  return (
    <CartContext.Provider value={{ cartItems, addToCart, removeFromCart }}>
      {children}
    </CartContext.Provider>
  )
}

export function useCart(): CartContextValue {
  const context = useContext(CartContext)
  if (!context) {
    throw new Error('useCart must be used within CartProvider')
  }
  return context
}
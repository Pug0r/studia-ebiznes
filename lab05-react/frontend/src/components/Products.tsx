import { useEffect, useState } from 'react'
import axios from 'axios'
import { API_URL } from '../config'
import { useCart } from '../hooks/useCart'
import type { Product } from '../types'

export default function Products() {
  const { addToCart } = useCart()
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function loadProducts(): Promise<void> {
      try {
        const response = await axios.get<Product[]>(`${API_URL}/products`)
        setProducts(response.data)
      } catch {
        setError('Nie udało się pobrać produktów')
      } finally {
        setLoading(false)
      }
    }

    loadProducts()
  }, [])

  return (
    <section>
      <h2>Produkty</h2>
      {loading && <p>Ładowanie...</p>}
      {error && <p>{error}</p>}
      {!loading && !error && (
        <ul>
          {products.map((product) => (
            <li key={product.id}>
              {product.name} - {product.price} PLN
              <button type="button" onClick={() => addToCart(product)}>
                Dodaj do koszyka
              </button>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
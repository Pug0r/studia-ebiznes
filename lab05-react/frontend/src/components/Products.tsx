import { useEffect, useState } from 'react'

type Product = {
  id: number
  name: string
  price: number
}

type ProductsProps = {
  apiUrl: string
}

export default function Products({ apiUrl }: ProductsProps) {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function loadProducts(): Promise<void> {
      try {
        const response = await fetch(`${apiUrl}/products`)
        if (!response.ok) {
          throw new Error('Błąd pobierania produktów')
        }
        const data: Product[] = await response.json()
        setProducts(data)
      } catch {
        setError('Nie udało się pobrać produktów')
      } finally {
        setLoading(false)
      }
    }

    loadProducts()
  }, [apiUrl])

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
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
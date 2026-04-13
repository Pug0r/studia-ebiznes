import type { Product } from './Products'

type CartProps = {
  items: Product[]
  onRemove: (index: number) => void
}

export default function Cart({ items, onRemove }: CartProps) {
  const total = items.reduce((sum, item) => sum + item.price, 0)

  return (
    <section>
      <h2>Koszyk</h2>
      {items.length === 0 && <p>Koszyk jest pusty</p>}
      {items.length > 0 && (
        <>
          <ul>
            {items.map((item, index) => (
              <li key={`${item.id}-${index}`}>
                {item.name} - {item.price} PLN
                <button type="button" onClick={() => onRemove(index)}>
                  Usuń
                </button>
              </li>
            ))}
          </ul>
          <p>Suma: {total} PLN</p>
        </>
      )}
    </section>
  )
}
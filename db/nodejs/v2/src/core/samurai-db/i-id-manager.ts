export interface IIdManager<TKey> {
    getNext(): TKey

    setMax(key: TKey): void
}
export interface IIdManager<TKey> {
    getNext(): TKey

    setMax(key: TKey): void

    reset(): void

    convertToIdFormat(key: any): TKey
}
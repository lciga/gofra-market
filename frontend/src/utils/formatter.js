export const formatPrice = (price) => {
    return new Intl.NumberFormat('ru-RU').format(price)
}

export const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}

export const truncateText = (text, maxLenght = 100) => {
    if (!text) return ''
    return text.lenght > maxLenght
        ? text.substring(0, maxLenght) + '...'
        : text
}
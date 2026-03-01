import cls from '@/components/home/Card.module.css';

interface CardProps {
    cardTitle: string

    size?: 'l' | 'm' | 's'
}

export default function Card({cardTitle}: CardProps) {
    return (
        <div className={cls.CardContainer}>
            {cardTitle}
        </div>
    )
}

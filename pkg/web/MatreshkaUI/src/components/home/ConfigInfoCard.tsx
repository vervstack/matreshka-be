import cls from '@/components/home/Card.module.css';

interface ConfigInfoCardProps {
    cardTitle: string

    size?: 'l' | 'm' | 's'
}

export default function ConfigInfoCard({cardTitle}: ConfigInfoCardProps) {
    return (
        <div className={cls.CardContainer}>
            {cardTitle}
        </div>
    )
}

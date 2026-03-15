import { useNavigate } from "react-router-dom";
import cls from '@/components/home/Card.module.css';

interface ConfigInfoCardProps {
    cardTitle: string

    size?: 'l' | 'm' | 's'
}

export default function ConfigInfoCard({cardTitle}: ConfigInfoCardProps) {
    const navigate = useNavigate();

    return (
        <div 
            className={cls.CardContainer}
            onClick={() => navigate(`/${cardTitle}`)}
            style={{ cursor: 'pointer' }}
        >
            {cardTitle}
        </div>
    )
}

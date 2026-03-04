import cn from "classnames";
import cls from "@/components/shared/SelectOneOf.module.css";

interface SelectOneOfProps {
    label?: string
    options: string[]
    selectedOption: string
    onSelected?: (option: string) => void
}

export default function SelectOneOf({label, options, selectedOption, onSelected}: SelectOneOfProps) {
    return (
        <div className={cls.SelectOneOfContainer}>
            {label && <label className={cls.Label}>{label}</label>}
            <div className={cls.OptionsWrapper}>
                {options.map((option) => (
                    <div
                        key={option}
                        className={cn(cls.Option, {
                            [cls.Selected]: option === selectedOption,
                        })}
                        onClick={() => onSelected?.(option)}
                    >
                        {option}
                    </div>
                ))}
            </div>
        </div>
    )
}

import {useState, useRef, useEffect} from "react";
import cn from "classnames";
import cls from "@/components/shared/SelectDropDown.module.css";
import ChevronDown from "@/assets/icons/ChevronDown.svg";

interface SelectDropDownProps {
    label?: string;
    options: string[];
    selectedOption: string;
    onSelected?: (option: string) => void;
}

export default function SelectDropDown({label, options, selectedOption, onSelected}: SelectDropDownProps) {
    const [isOpen, setIsOpen] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        function handleClickOutside(event: MouseEvent) {
            if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const toggleOpen = () => setIsOpen(!isOpen);

    const handleSelect = (option: string) => {
        onSelected?.(option);
        setIsOpen(false);
    };

    return (
        <div className={cls.SelectDropDownContainer} ref={containerRef}>
            <div className={cn(cls.SelectedDisplay, {[cls.Open]: isOpen})} onClick={toggleOpen}>
                <span className={cls.SelectedValue}>{selectedOption}</span>
                <img src={ChevronDown} className={cn(cls.Chevron, {[cls.Open]: isOpen})} alt="▼" />
                {label && (
                    <label className={cn(cls.Label, cls.Floating)}>
                        {label}
                    </label>
                )}
            </div>
            {isOpen && (
                <div className={cls.OptionsWrapper}>
                    {options.map((option) => (
                        <div
                            key={option}
                            className={cn(cls.Option, {
                                [cls.Selected]: option === selectedOption,
                            })}
                            onClick={() => handleSelect(option)}
                        >
                            {option}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

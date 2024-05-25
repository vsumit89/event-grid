export function Backdrop({ children, onClose }) {
    return (
        <div
            className="fixed inset-0 z-10 bg-black bg-opacity-50 w-screen h-screen flex justify-center items-center"
            onClick={onClose}
        >
            {children}
        </div>
    )
}

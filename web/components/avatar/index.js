export const Avatar = ({ image, name }) => {
    return (
        <div
            className={`w-10 h-10 rounded-full bg-secondary-background flex items-center justify-center text-white ${
                image ? 'overflow-hidden' : ''
            }`}
        >
            {image ? (
                <img
                    src={image}
                    alt="Avatar"
                    className="w-full h-full object-cover"
                />
            ) : (
                name.charAt(0).toUpperCase()
            )}
        </div>
    )
}

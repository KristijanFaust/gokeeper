import './button.styles.scss';

const Button = ({children, disabled, onClick, tooltip}) => (
  <button title={tooltip ?? null} className={disabled ? 'button disabled' : 'button'} disabled={disabled} onClick={onClick}>
    {children}
  </button>
);

export default Button;

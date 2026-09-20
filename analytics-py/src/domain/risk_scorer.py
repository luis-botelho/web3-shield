def calculate_risk_score(signature: str) -> int:
    """Retorna a classificação de risco associada ao seletor ABI informado."""
    
    # approve pode conceder autorização de gasto a contratos de terceiros.
    if signature == "0x095ea7b3":
        return 90
        
    # transfer representa uma transferência ERC-20 convencional.
    elif signature == "0xa9059cbb":
        return 30
        
    # Seletores desconhecidos e transferências nativas usam o risco base.
    return 10

# risk_scorer.py

def calculate_risk_score(signature: str) -> int:
    """
    Recebe a assinatura da função (ex: 0x095ea7b3) e retorna um score de risco de 0 a 100.
    """
    
    # 0x095ea7b3 = Função 'approve' (Permite que outro contrato gaste seus tokens)
    # Risco altíssimo, muito usada em phishing e drains.
    if signature == "0x095ea7b3":
        return 90
        
    # 0xa9059cbb = Função 'transfer' (Transferência de tokens ERC-20)
    # Risco moderado, transação comum.
    elif signature == "0xa9059cbb":
        return 30
        
    # Outras transações, contratos desconhecidos ou transferências nativas
    return 10
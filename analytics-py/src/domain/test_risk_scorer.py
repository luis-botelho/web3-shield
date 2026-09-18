from src.domain.risk_scorer import calculate_risk_score

def test_calculate_risk_score_approve_perigoso():
    # 0x095ea7b3 é a assinatura mundial padrão para a função "Approve" (Aprovar gastos)
    # É a função mais usada em golpes de drenagem de carteira. Risco Altíssimo!
    score = calculate_risk_score("0x095ea7b3")
    assert score == 90

def test_calculate_risk_score_transferencia_comum():
    # 0xa9059cbb é a assinatura da função "Transfer"
    score = calculate_risk_score("0xa9059cbb")
    assert score == 30

def test_calculate_risk_score_desconhecido():
    # Uma transferência de ETH nativo sem função ou algo que não conhecemos
    score = calculate_risk_score("0x00000000")
    assert score == 10
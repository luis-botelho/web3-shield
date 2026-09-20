from src.domain.risk_scorer import calculate_risk_score

def test_calculate_risk_score_approve_perigoso():
    score = calculate_risk_score("0x095ea7b3")
    assert score == 90

def test_calculate_risk_score_transferencia_comum():
    score = calculate_risk_score("0xa9059cbb")
    assert score == 30

def test_calculate_risk_score_desconhecido():
    score = calculate_risk_score("0x00000000")
    assert score == 10

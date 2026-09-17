import time
import sys

# Importamos a nossa Mão (Repositório) e o nosso Cérebro (Domínio)
from src.repository.postgres_repo import PostgresRepository
from src.domain.risk_scorer import calculate_risk_score

def main():
    print("🕵️‍♂️ Iniciando o Detetive (Web3 Shield Analytics)...")
    
    try:
        repo = PostgresRepository()
        print("🐘 Conectado ao banco de dados com sucesso!")
    except Exception as e:
        print(f"🚨 Erro fatal ao conectar no banco: {e}")
        sys.exit(1)

    print("🔄 Vigiando o banco de dados em busca de novas transações...")

    try:
        while True:
            # 1. Busca transações que o Go salvou, mas que o Python ainda não viu
            unprocessed = repo.get_unprocessed_transactions()
            
            if unprocessed:
                print(f"\n📦 Encontradas {len(unprocessed)} transações pendentes de análise.")
                
                for tx in unprocessed:
                    tx_hash = tx['tx_hash']
                    signature = tx['function_signature']
                    
                    # 2. Passa a assinatura na nossa regra de negócio (TDD validou isso!)
                    score = calculate_risk_score(signature)
                    
                    # 3. Salva o veredito no banco
                    repo.save_risk_analysis(tx_hash, score)
                    
                    print(f"  🔍 Analisado: {tx_hash[:10]}... | Sig: {signature} -> Nível de Risco: {score}")
            
            # Pausa por 3 segundos para não sobrecarregar o processador e o banco
            time.sleep(3)
            
    except KeyboardInterrupt:
        # Se você apertar Ctrl+C, ele fecha a conexão com o banco educadamente
        print("\n🛑 Encerrando o Detetive. Fechando conexão com o banco...")
    finally:
        repo.close()

if __name__ == "__main__":
    main()
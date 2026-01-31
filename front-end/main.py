import streamlit as st
import requests

st.set_page_config(layout="wide", page_title="Gopher gallery")

@st.cache_data(ttl=10)
def fetch_data():
    try:
        response = requests.get("http://localhost:8080/api/v1/files", timeout=5)
        if response.status_code == 200:
            return response.json()
    except Exception as e:
        st.error(f"Erro ao conectar na API: {e}")
    return []

def sendData():
    st.sidebar.header("Enviar nova imagem")
    imagem = st.sidebar.file_uploader(
        "Escolher imagem",
        type=["png", "jpg", "jpeg", "gif"],
        accept_multiple_files=False
    )

    if imagem is not None:
        if st.sidebar.button("Enviar imagem"):
            fe = {"file": (imagem.name, imagem.getvalue(), imagem.type)}
            try:
                with st.spinner("Enviando..."):
                    # Verifique se o endpoint Upfile está correto (case sensitive)
                    response = requests.post("http://localhost:8080/api/v1/Upfile", files=fe)

                if response.status_code in [200, 201]:
                    st.sidebar.success("Concluído!")
                    st.cache_data.clear()
                    st.rerun()
                else:
                    st.sidebar.error(f"Erro no upload: {response.status_code} - {response.text}")
            except Exception as e:
                st.sidebar.error(f"Erro de conexão: {e}")
                
def delete_image(image):
    try: 
        # CORREÇÃO: Adicionado http://
        url = f"http://localhost:8080/api/v1/Delefile/{image['fileKey']}"
        response = requests.delete(url)
        if response.status_code == 200:
            st.toast(f"{image['name']} deletado!")
            st.cache_data.clear()
            st.rerun()
        else:
            st.error(f"Erro ao deletar: {response.status_code}")
    except Exception as e:
        st.error(f"Erro: {e}")
            
@st.dialog("Zoom e opções")
def full_screen(image):
    # CORREÇÃO: Usar a URL para mostrar a imagem, não o nome
    st.image(image["url"])
    st.write(f"**Nome:** {image['name']}")
    st.write(f"**Key:** {image['fileKey']}")

    st.divider()

    cols = st.columns([1, 1])
    with cols[1]:
        if st.button("🗑️ Deletar Imagem", key=f"del_{image['fileKey']}", type="primary"):
            delete_image(image)
    with cols[0]:
        if st.button("Fechar", key=f"close_{image['fileKey']}"):
            st.rerun()

@st.fragment(run_every=10)
def render_gallery():
    data = fetch_data()

    if not data:
        st.info("Nenhuma imagem encontrada ou API offline.")
        return

    st.write(f"Galeria atualizada: {len(data)} imagens encontradas.")

    cols_count = 4
    for i in range(0, len(data), cols_count):
        row = data[i : i + cols_count]
        colunas = st.columns(cols_count)
        for j, image_data in enumerate(row):
            with colunas[j]:
                st.image(
                    image_data['url'],
                    caption=image_data['name'],
                )
                if st.button("Abrir", key=f"btn_view_{image_data['fileKey']}", use_container_width=True):
                    full_screen(image_data)

# --- Execução Principal ---
st.title("Gopher Gallery")
st.markdown("---")

sendData()      # Sidebar primeiro
render_gallery() # Galeria depois
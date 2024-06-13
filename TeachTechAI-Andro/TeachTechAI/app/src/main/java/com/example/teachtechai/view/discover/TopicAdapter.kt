package com.example.teachtechai.view.discover

import android.content.Intent
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.example.teachtechai.R
import com.example.teachtechai.databinding.DashboardCardviewBinding
import com.example.teachtechai.view.chat.ChatActivity

data class Topic(val title: String, val description: String, val imageResId: Int)


class TopicAdapter(private val topics: List<Topic>) : RecyclerView.Adapter<TopicAdapter.TopicViewHolder>() {

    class TopicViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        val topicImage: ImageView = itemView.findViewById(R.id.topic_image)
        val topicTitle: TextView = itemView.findViewById(R.id.topic_title)
        val topicDescription: TextView = itemView.findViewById(R.id.topic_description)
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): TopicViewHolder {
        val view = LayoutInflater.from(parent.context).inflate(R.layout.dashboard_cardview, parent, false)
        return TopicViewHolder(view)
    }

    override fun onBindViewHolder(holder: TopicViewHolder, position: Int) {
        val topic = topics[position]
        holder.topicImage.setImageResource(topic.imageResId)
        holder.topicTitle.text = topic.title
        holder.topicDescription.text = topic.description

        holder.itemView.setOnClickListener{
            val intent = Intent(holder.itemView.context, ChatActivity::class.java)
            intent.putExtra("title", topic.title)
            holder.itemView.context.startActivity(intent)
        }
    }
    override fun getItemCount(): Int = topics.size
}